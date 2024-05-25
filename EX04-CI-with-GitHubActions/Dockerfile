FROM golang:1.22-alpine

# Set the maintainer of the image
LABEL at.hagenberg.mc.image.authors="Stefan.Penzinger@fh-hagenberg.at"

# Set working directory: `/src`
WORKDIR /src

# Copy local file `main.go` to the working directory
COPY main.go /src
COPY go.mod /src

# List items in the working directory (ls)
RUN ls

# Build the GO app as myapp binary and move it to /usr/
RUN go build -o /usr/myapp

#Expose port 8888
EXPOSE 8888

# Run the service myapp when a container of this image is launched
CMD ["/usr/myapp"]
