# Robot Warehouse Task Solution

## Demo link

[Please check the demo](https://dronedeploy-challenge.s3.ap-southeast-2.amazonaws.com/demo.mov)

## To run the app

1. Start local backend first
   **required go version 1.25**
   ```sh
   cd backend && go run main.go
   ```
2. Start local frontend
   **required node version 20**
   ```sh
   cd frontend
   npm install && npm run dev
   ```

## Solution overview:

I chose a single-flight queue design — meaning only one task can be queued at a time.

The key concern is boundary safety:

If multiple tasks are queued, each task wouldn’t know the final position of the robot after the previous task as robot may fail midway.

#### Example:

Task 1 → NN, Task 2 → SS.

If everything goes well, the robot moves back to the original point.

But if the robot fails after completing the first N, then the two S commands from Task 2 will move the robot to (0, -1), which is outside the boundary.

We only know the robot’s final position when Task 1 has reached a terminal state: COMPLETED, CANCELLED, or FAILED.

To avoid this, a new task is only accepted when the previous task has reached a terminal state.

This ensures each task is executed safely without conflicting robot positions.

---

Pretty enojoy my 2 day GO learning here.

## Requirements

[x] Create a RESTful API to accept a series of commands to the robot.

[x] Make sure that the robot doesn't try to move outside the warehouse.

[x] Create a RESTful API to report the command series's execution status.

[x] Create a RESTful API cancel the command series.

[x] The RESTful service should be written in Golang.

## Challenge

[x] The Robot SDK is still under development, you need to find a way to prove your API logic is working.

**Answer**

1. Integration tests – Assume the SDK always returns some data, and verify that the API behaves as expected.

2. Frontend demo – Demonstrates that the API is working correctly and that errors are handled gracefully.

---

[x] The ground control station wants to be notified as soon as the command sequence completed. Please provide a high level design overview how you can achieve it. This overview is not expected to be hugely detailed but should clearly articulate the fundamental concept in your design.

**Answer**

Event-driven architecture

- [Please View the screenshot from S3](https://dronedeploy-challenge.s3.ap-southeast-2.amazonaws.com/notification-architecture.png).

In this app, task information is stored in memory with a status of PENDING, COMPLETED, FAILED, or CANCELLED.

In production, this would be backed by a proper database. A database record update event can then be emitted whenever a task’s status changes.

lets use aws service as an example.

say all of the the task information saved in the database, if the existing record update in the database changes, then it will trigger the
Lambda function.

A Lambda function can listen for these events.
If the task status is COMPLETED, the Lambda publishes an event to an SNS topic.

The Ground Control Station subscribes to the SNS topic. so they can be notified via email, message, or even pager duty.
