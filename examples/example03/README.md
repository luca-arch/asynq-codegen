# Example 03

This directory contains a tiny, self-contained, working example. 

It uses Docker Compose to spin up a Redis container, along one consumer and one producer containers.

Open a terminal in this folder and run

```sh
docker compose up
```

Once the build finishes, you should see a never ending output that looks like

```
producer-1  | 2026/01/26 01:53:53 INFO new task pushed task.id=25933200-0800-4240-be05-a1e23d191b6b task.type=main:task_2 type=main:task_2
consumer-1  | 2026/01/26 01:53:53 INFO processing task created="2026-01-26 01:53:53.522479761 +0000 UTC" type=main:task_2
producer-1  | 2026/01/26 01:53:54 INFO new task pushed task.id=dd17895b-dc38-4c3a-a891-076eb8c29e6d task.type=main:task_1 type=main:task_1
producer-1  | 2026/01/26 01:53:55 INFO new task pushed task.id=780330b0-ef90-41f1-83f5-320bc963323d task.type=main:task_2 type=main:task_2
consumer-1  | 2026/01/26 01:53:55 INFO processing task created="2026-01-26 01:53:54.523509678 +0000 UTC" type=main:task_1
consumer-1  | 2026/01/26 01:53:55 INFO processing task created="2026-01-26 01:53:55.524719804 +0000 UTC" type=main:task_2
```

Press `CTRL + C` (or `CMD + C` on Mac) to terminate the example, then run `docker compose down -v` to tear down the containers.
