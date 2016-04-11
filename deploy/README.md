
### Contiv Deploy

Deploy utilizes libcompose to run applications and apply policies automatically to them. At this point it applies security policies between application groups

#### How to try it out

1. Checkout the tree: 
git clone https://github.com/jainvipin/libcompose

2. Compile and run unit tests to ensure you have correct environment
cd libcompose
make binary

3. Launch containers
```
$ cd example
$ deploy -file docker-compose.yml --labels="io.contiv.env:prod"
```

4. Verify that policies are instantiated


5. Use docker-compose as usual to ps/stop/scale containers
```
$ docker-compose ps
$ docker-compose stop
$ docker-compose restart
```
