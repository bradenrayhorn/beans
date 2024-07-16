FROM node:20-alpine@sha256:34b7aa411056c85dbf71d240d26516949b3f72b318d796c26b57caaa1df5639a as base

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
