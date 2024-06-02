FROM node:20-alpine@sha256:928b24aaadbd47c1a7722c563b471195ce54788bf8230ce807e1dd500aec0549 as base

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
