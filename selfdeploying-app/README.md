## Selfdeploying App Example

This example is inspired by a talk from Kelsey Hightower: https://www.youtube.com/watch?v=XPC-hFL-4lU
The idea is that if you need to deploy stuff, why should you learn how to deploy my app to kubernetes? Why shouldn't it be able to do so by itself? And this one can! 

### Usage

Run locally on port 8000.
```
    dist/app
```

Deploy to kubernetes, still port 8000.
```
    dist/app --kubernetes --replicas 2
```