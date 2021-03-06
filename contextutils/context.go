// Contains common flyte context utils.
package contextutils

import (
	"context"
	"fmt"
)

type Key string

const (
	AppNameKey    Key = "app_name"
	NamespaceKey  Key = "ns"
	TaskTypeKey   Key = "tasktype"
	ProjectKey    Key = "project"
	DomainKey     Key = "domain"
	WorkflowIDKey Key = "wf"
	NodeIDKey     Key = "node"
	TaskIDKey     Key = "task"
	ExecIDKey     Key = "exec_id"
	JobIDKey      Key = "job_id"
	PhaseKey      Key = "phase"
)

func (k Key) String() string {
	return string(k)
}

var logKeys = []Key{
	AppNameKey,
	JobIDKey,
	NamespaceKey,
	ExecIDKey,
	NodeIDKey,
	WorkflowIDKey,
	TaskTypeKey,
	PhaseKey,
}

// Gets a new context with namespace set.
func WithNamespace(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, NamespaceKey, namespace)
}

// Gets a new context with JobId set. If the existing context already has a job id, the new context will have
// <old_jobID>/<new_jobID> set as the job id.
func WithJobID(ctx context.Context, jobID string) context.Context {
	existingJobID := ctx.Value(JobIDKey)
	if existingJobID != nil {
		jobID = fmt.Sprintf("%v/%v", existingJobID, jobID)
	}

	return context.WithValue(ctx, JobIDKey, jobID)
}

// Gets a new context with AppName set.
func WithAppName(ctx context.Context, appName string) context.Context {
	return context.WithValue(ctx, AppNameKey, appName)
}

// Gets a new context with Phase set.
func WithPhase(ctx context.Context, phase string) context.Context {
	return context.WithValue(ctx, PhaseKey, phase)
}

// Gets a new context with ExecutionID set.
func WithExecutionID(ctx context.Context, execID string) context.Context {
	return context.WithValue(ctx, ExecIDKey, execID)
}

// Gets a new context with NodeID (nested) set.
func WithNodeID(ctx context.Context, nodeID string) context.Context {
	existingNodeID := ctx.Value(NodeIDKey)
	if existingNodeID != nil {
		nodeID = fmt.Sprintf("%v/%v", existingNodeID, nodeID)
	}
	return context.WithValue(ctx, NodeIDKey, nodeID)
}

// Gets a new context with WorkflowName set.
func WithWorkflowID(ctx context.Context, workflow string) context.Context {
	return context.WithValue(ctx, WorkflowIDKey, workflow)
}

// Get new context with Project and Domain values set
func WithProjectDomain(ctx context.Context, project, domain string) context.Context {
	c := context.WithValue(ctx, ProjectKey, project)
	return context.WithValue(c, DomainKey, domain)
}

// Gets a new context with WorkflowName set.
func WithTaskID(ctx context.Context, taskID string) context.Context {
	return context.WithValue(ctx, TaskIDKey, taskID)
}

// Gets a new context with TaskType set.
func WithTaskType(ctx context.Context, taskType string) context.Context {
	return context.WithValue(ctx, TaskTypeKey, taskType)
}

func addFieldIfNotNil(ctx context.Context, m map[string]interface{}, fieldKey Key) {
	val := ctx.Value(fieldKey)
	if val != nil {
		m[fieldKey.String()] = val
	}
}

func addStringFieldWithDefaults(ctx context.Context, m map[string]string, fieldKey Key) {
	val := ctx.Value(fieldKey)
	if val == nil {
		val = ""
	}
	m[fieldKey.String()] = val.(string)
}

// Gets a map of all known logKeys set on the context. logKeys are special and should be used incase, context fields
// are to be added to the log lines.
func GetLogFields(ctx context.Context) map[string]interface{} {
	res := map[string]interface{}{}
	for _, k := range logKeys {
		addFieldIfNotNil(ctx, res, k)
	}
	return res
}

func Value(ctx context.Context, key Key) string {
	val := ctx.Value(key)
	if val != nil {
		return val.(string)
	}
	return ""
}

func Values(ctx context.Context, keys ...Key) map[string]string {
	res := map[string]string{}
	for _, k := range keys {
		addStringFieldWithDefaults(ctx, res, k)
	}
	return res
}
