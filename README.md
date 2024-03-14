# Cronih
## scheduling hit your API with Basic Auth


### Get Started

1. clone the repository
1. rename jobs.yaml.example to jobs.yaml
1. build and running


### Usage

example job with schedule 
```
jobs:
  job1minutes:
    enable: true
    schedule: "* * * * *"
    periodic: 0s
    endpoint: "http://localhost:8080/job1minutes"
    credentials:
      username: username
      password: password
```


example job with periodic 
```
jobs:
  job5s:
    enable: true
    schedule: ""
    periodic: 5s
    endpoint: "http://localhost:8080/job5s"
    credentials:
      username: username
      password: password
```

example job with periodic 
```
jobs:
  jobDisableJob:
    enable: false
    schedule: ""
    periodic: 5s
    endpoint: "http://localhost:8080/jobDisableJob"
    credentials:
      username: username
      password: password
```