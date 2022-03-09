UNAME := $(shell uname -s)
SRCDIR	 = src
PRVDIR	 = priv

ifeq ($(UNAME),Linux)
TARGET = $(PRVDIR)/sport_linux
endif

ifeq ($(UNAME),Darwin)
TARGET = $(PRVDIR)/sport_darwin
endif

SOURCES = ${SRCDIR}/*.go

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(SOURCES)
	mkdir -p $(PRVDIR)
	(cd $(SRCDIR); go build -o ../$(TARGET))

clean:
	rm -f $(TARGET)
