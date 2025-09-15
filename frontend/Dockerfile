# Use an official Node.js runtime as a parent image
FROM node:20-alpine

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json (if available)
COPY package*.json ./

# Install project dependencies
RUN npm install

# Copy the rest of the application's source code
# This is useful for building the image, though the docker-compose volume will override it
COPY . .

# The React development server runs on port 3000
EXPOSE 3000

# The command to start the development server
CMD ["npm", "start"]
