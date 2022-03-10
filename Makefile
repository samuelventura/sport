UNAME := $(shell uname -s)
SRCDIR	 = src
PRVDIR	 = priv
TARGET = $(PRVDIR)/sport

ifeq ($(UNAME),Linux)
endif

ifeq ($(UNAME),Darwin)
endif

SOURCES = ${SRCDIR}/*.go

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(SOURCES)
	mkdir -p $(PRVDIR)
	(cd $(SRCDIR); go build -o ../$(TARGET))

clean:
	rm -f $(TARGET)
