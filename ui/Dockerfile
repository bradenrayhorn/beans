FROM node:20-alpine@sha256:c13b26e7e602ef2f1074aef304ce6e9b7dd284c419b35d89fcf3cc8e44a8def9 as base

# build step
FROM base as build

RUN mkdir /app
COPY / app/
WORKDIR /app
RUN npm install
RUN npm run build

# final image
FROM base

RUN mkdir -p /app/build
COPY --from=build /app/build /app/package.json /app/entry.js /app/

ENV NODE_ENV=production PORT=8080

ARG GIT_SHA
ENV PUBLIC_VERSION=${GIT_SHA}

CMD ["node", "/app/entry.js"]
