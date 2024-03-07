FROM ubuntu:latest as build

# Install dependencies
RUN apt-get update && \
    apt-get install -y \
    build-essential \
    autoconf \
    automake \
    libtool \
    pkg-config \
    wget \
    bzip2

# Get hwloc source code
WORKDIR /tmp
ENV HWLOC_MAJOR_VERSION=2.10
ENV HWLOC_MINOR_VERSION=0
ENV HWLOC_VERSION=${HWLOC_MAJOR_VERSION}.${HWLOC_MINOR_VERSION}
RUN wget https://download.open-mpi.org/release/hwloc/v${HWLOC_MAJOR_VERSION}/hwloc-${HWLOC_VERSION}.tar.bz2
RUN tar -xjf hwloc-${HWLOC_VERSION}.tar.bz2

# Build hwloc
WORKDIR /tmp/hwloc-${HWLOC_VERSION}
RUN ./configure && \
    make && \
    make install && \
    make install DESTDIR=/tmp/hwloc

FROM golang:1.21.7-bookworm

# Copy hwloc binaries and libraries from the builder stage
COPY --from=build /tmp/hwloc/usr/local/lib/ /usr/local/lib/
COPY --from=build /tmp/hwloc/usr/local/include/ /usr/local/include/

# Configure env values
ENV C_INCLUDE_PATH /usr/local/include:$C_INCLUDE_PATH
ENV LD_LIBRARY_PATH usr/local/lib:$LD_LIBRARY_PATH

WORKDIR $GOPATH
