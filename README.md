# Golang Rest Api

## Techs
* Redis
* Nginx
* Docker/docker-compose

## Start Server
```
cd golang-rest-api
sudo docker-compose up
```

## APIS

### http://localhost/records  (AWS URL: ec2-3-121-109-138.eu-central-1.compute.amazonaws.com/records)
This end point will return the query results from Getir mongoDB collection.

#### POST
Only post method is allowed.  

* Sample Request Payload:  
```javascript
{
  "startDate": "2016-01-01",
  "endDate": "2016-02-02",
  "minCount": 2700,
  "maxCount": 2800
}
```

* Sample Response Payload:
```javascript
{
    "code": 0,
    "msg": "Success",
    "records": [
        {
            "key": "xdekszeS",
            "createdAt": "2016-01-25T15:25:24Z0",
            "totalCount": 2739
        },
        {
            "key": "iNAwtoZQ",
            "createdAt": "2016-01-22T03:22:07Z0",
            "totalCount": 2705
        }
    ]
}
```

### http://localhost/in-memory  (AWS URL: ec2-3-121-109-138.eu-central-1.compute.amazonaws.com/in-memory)

#### POST
You can send post request to the api to save a key value pair (both string) to the redis database.

* Sample Request Payload:  
```javascript 
{  
  "key": "active-tabs",  
  "value": "getir"  
}  
```
* Response Payload:  
```javascript 
{  
  "key": "active-tabs",  
  "value": "getir"  
}  
```

#### GET  
You can send get request to /in-memory with a key parameter in url.  
```
http://localhost/in-memory?key=active-tabs
```
* Response Payload:  
```javascript 
{  
  "key": "active-tabs",  
  "value": "getir"  
}  
```


