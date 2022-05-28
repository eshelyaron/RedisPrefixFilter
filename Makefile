all:
	g++ -shared --std=c++17 -march=native -fPIC -o module.so module.c Prefix-Filter/Tests/*.cpp Prefix-Filter/*.cpp Prefix-Filter/Prefix-Filter/*.cpp Prefix-Filter/TC-Shortcut/*.cpp
