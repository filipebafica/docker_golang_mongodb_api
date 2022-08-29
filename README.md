# docker_golang_mongodb_api
This is a REST api with docker + golang + mongodb

## ⚙️ Getting Started
```
$ git clone https://github.com/filipebafica/docker_golang_mongodb_api.git
$ cd docker_golang_mongodb_api
$ docker compose build
$ docker compose up -d
```

## 🌐 How to Use
Requests may be done to the folloing end points:

`[POST] 0.0.0.0:8000/person` \
`[GET] 0.0.0.0:8000/people`

Available fields:
```
  {
    "firstname": "Eddie",
    "lastname": "Head",
  }

```
