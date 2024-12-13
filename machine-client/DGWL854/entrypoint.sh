#!/bin/sh
exec service ssh start
exec /main $@
