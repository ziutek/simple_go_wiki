include $(GOROOT)/src/Make.inc

TARG=wiki
GOFILES=\
	view.go\
	controller.go\
	mysql.go\

include $(GOROOT)/src/Make.cmd
