# Stage: UI
FROM node:12.18-alpine as ui

# Create app directory
WORKDIR /app

# Bundle app source
COPY . /app/

# Disable Source Map
ENV GENERATE_SOURCEMAP=false

# Install packages
RUN npm install
RUN npm run build


# Stage: Core
FROM nginx:1.19-alpine

# Copy ui
COPY --from=ui /app/build/ /usr/share/nginx/html