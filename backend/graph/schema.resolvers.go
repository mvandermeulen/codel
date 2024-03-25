package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/semanser/ai-coder/database"
	"github.com/semanser/ai-coder/executor"
	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/graph/subscriptions"
)

// CreateFlow is the resolver for the createFlow field.
func (r *mutationResolver) CreateFlow(ctx context.Context) (*gmodel.Flow, error) {
	flow, err := r.Db.CreateFlow(ctx, database.CreateFlowParams{
		Name:   database.StringToPgText("New Task"),
		Status: database.StringToPgText("in_progress"),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create flow: %w", err)
	}

	return &gmodel.Flow{
		ID:     uint(flow.ID),
		Name:   flow.Name.String,
		Status: gmodel.FlowStatus(flow.Status.String),
	}, nil
}

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, flowID uint, query string) (*gmodel.Task, error) {
	type InputTaskArgs struct {
		Query string `json:"query"`
	}

	args := InputTaskArgs{Query: query}
	arg, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	task, err := r.Db.CreateTask(ctx, database.CreateTaskParams{
		Type:    database.StringToPgText("input"),
		Message: database.StringToPgText(query),
		Status:  database.StringToPgText("finished"),
		Args:    arg,
		FlowID:  pgtype.Int8{Int64: int64(flowID), Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	executor.AddCommand(task)

	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	return &gmodel.Task{
		ID:        uint(task.ID),
		Message:   task.Message.String,
		Type:      gmodel.TaskType(task.Type.String),
		Status:    gmodel.TaskStatus(task.Status.String),
		Args:      string(task.Args),
		CreatedAt: task.CreatedAt.Time,
	}, nil
}

// Exec is the resolver for the _exec field.
func (r *mutationResolver) Exec(ctx context.Context, containerID string, command string) (string, error) {
	b := bytes.Buffer{}
	executor.ExecCommand(containerID, command, &b)

	return b.String(), nil
}

// Flows is the resolver for the flows field.
func (r *queryResolver) Flows(ctx context.Context) ([]*gmodel.Flow, error) {
	flows, err := r.Db.ReadAllFlows(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch flows: %w", err)
	}

	var gFlows []*gmodel.Flow

	for _, flow := range flows {
		var gTasks []*gmodel.Task

		gFlows = append(gFlows, &gmodel.Flow{
			ID:            uint(flow.ID),
			Name:          flow.Name.String,
      Terminal: &gmodel.Terminal{
        ContainerName: flow.ContainerName.String,
        Available: false,
      },
			Tasks:         gTasks,
			Status:        gmodel.FlowStatus(flow.Status.String),
		})
	}

	return gFlows, nil
}

// Flow is the resolver for the flow field.
func (r *queryResolver) Flow(ctx context.Context, id uint) (*gmodel.Flow, error) {
	flow, err := r.Db.ReadFlow(ctx, int64(id))

	if err != nil {
		return nil, fmt.Errorf("failed to fetch flow: %w", err)
	}

	var gFlow *gmodel.Flow
	var gTasks []*gmodel.Task

	tasks, err := r.Db.ReadTasksByFlowId(ctx, pgtype.Int8{Int64: int64(id), Valid: true})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	for _, task := range tasks {
		gTasks = append(gTasks, &gmodel.Task{
			ID:        uint(task.ID),
			Message:   task.Message.String,
			Type:      gmodel.TaskType(task.Type.String),
			Status:    gmodel.TaskStatus(task.Status.String),
			Args:      string(task.Args),
			Results:   task.Results.String,
			CreatedAt: task.CreatedAt.Time,
		})
	}

	gFlow = &gmodel.Flow{
		ID:            uint(flow.ID),
		Name:          flow.Name.String,
		Tasks:         gTasks,
    Terminal: &gmodel.Terminal{
      ContainerName: flow.ContainerName.String,
      Available: flow.ContainerStatus.String == "running",
    },
		Status:        gmodel.FlowStatus(flow.Status.String),
	}

	return gFlow, nil
}

// TaskAdded is the resolver for the taskAdded field.
func (r *subscriptionResolver) TaskAdded(ctx context.Context, flowID uint) (<-chan *gmodel.Task, error) {
	return subscriptions.TaskAdded(ctx, int64(flowID))
}

// TaskUpdated is the resolver for the taskUpdated field.
func (r *subscriptionResolver) TaskUpdated(ctx context.Context) (<-chan *gmodel.Task, error) {
	panic(fmt.Errorf("not implemented: TaskUpdated - taskUpdated"))
}

// FlowUpdated is the resolver for the flowUpdated field.
func (r *subscriptionResolver) FlowUpdated(ctx context.Context, flowID uint) (<-chan *gmodel.Flow, error) {
	return subscriptions.FlowUpdated(ctx, int64(flowID))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
