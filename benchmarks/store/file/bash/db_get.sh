#!/bin/bash

grep "^$1," database | sed -e "s/^$1,//" | tail -n 1
