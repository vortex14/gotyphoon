**Лдлдл**

1.

2.

![Снимок экрана 2021-09-25 в 09.51.39](/Users/vortex/Desktop/Снимок экрана 2021-09-25 в 09.51.39.png)3.



```bash
# create the volume in advance
  $ docker volume create --driver local \
      --opt type=none \
      --opt device=/home/user/test \
      --opt o=bind \
      test_vol

  # create on the fly with --mount
  $ docker run -it --rm \
    --mount type=volume,dst=/container/path,volume-driver=local,volume-opt=type=none,volume-opt=o=bind,volume-opt=device=/home/user/test \
    foo

  # inside a docker-compose file
  ...
  volumes:
    bind-test:
      driver: local
      driver_opts:
        type: none
        o: bind
        device: /home/user/test
  ...
```