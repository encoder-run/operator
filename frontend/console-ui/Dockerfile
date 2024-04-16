# Step 1: Build the React application
FROM node:18 AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . ./
RUN npm run build

# Step 2: Serve the application using a lightweight node image
FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/build /app
RUN npm install -g serve
CMD ["serve", "-s", ".", "-l", "8081"]

# Expose port 8081 to the outside once the container is running
EXPOSE 8081
