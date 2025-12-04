package ksuid

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/servicecontext"
)

// Production is the internal name for production ksuid, but is omitted
// during marshaling.
const Production = "prod"

var exportedNode = makeNode(context.Background(), Production)

func makeNode(ctx context.Context, environment string) *Node {
	if iid, err := NewDockerID(ctx); err == nil {
		return NewNode(environment, iid)
	}

	if iid, err := NewHardwareID(ctx); err == nil {
		return NewNode(environment, iid)
	}

	return NewNode(environment, NewRandomID())
}

// Node contains metadata used for ksuid generation for a specific machine.
type Node struct {
	InstanceID InstanceID

	timestamp  uint64
	sequence   uint32
	sequenceMu sync.Mutex
}

// NewNode returns a ID generator for the current machine.
func NewNode(environment string, instanceID InstanceID) *Node {
	return &Node{
		InstanceID: instanceID,
	}
}

// Generate returns a new ID for the machine and resource configured.
func (n *Node) Generate(ctx context.Context, resource string) (id ID) {
	if strings.ContainsRune(resource, '_') {
		panic(merr.New(ctx, "ksuid_resource_contains_underscore", merr.M{
			"resource": resource,
		}))
	}

	if info := servicecontext.GetContext(ctx); info != nil {
		id.Environment = info.Env
	} else {
		id.Environment = Production
	}

	id.Resource = resource
	id.InstanceID = n.InstanceID

	n.sequenceMu.Lock()

	timestamp := uint64(time.Now().UTC().Unix())
	if (timestamp - n.timestamp) >= 1 {
		n.timestamp = timestamp
		n.sequence = 0
	} else {
		n.sequence++
	}

	id.Timestamp = timestamp
	id.SequenceID = n.sequence

	n.sequenceMu.Unlock()

	return id
}

// SetInstanceID overrides the default instance id in the exported node.
// This will effect all invocations of the Generate function.
func SetInstanceID(instanceID InstanceID) {
	exportedNode.InstanceID = instanceID
}

// Generate returns a new ID for the current machine and resource configured.
func Generate(ctx context.Context, resource string) ID {
	return exportedNode.Generate(ctx, resource)
}
