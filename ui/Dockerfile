FROM node:20-alpine@sha256:df01469346db2bf1cfc1f7261aeab86b2960efa840fe2bd46d83ff339f463665 as base

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
