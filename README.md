# goTyphoon

GoTyphoon must provide a few optimized way between python and go version.


GoTyphoon can be replaced PyTyphoon, except processor component or only part of processor (pipeline).

Symbian mode will be pass task to python pipeline and return to golang env as finish operation.

Rust PyO3 provided better way for python extension better that glue of DataDog company. We need return after sometime for deep research

For Typhoon Micro build:
1. All component running in one process, may be in different thread

For Typhoon bundle build:
1. All projects of cluster will be running in one process

For Typhoon Symbian mode:
1. Pass Typhoon task to Python Env
2. Finish task only into golang Env
2. import python module for pipeline