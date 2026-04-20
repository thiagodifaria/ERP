# Build stage for the finance service.
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /src

COPY Directory.Build.props ./
COPY src/Finance.Api/Finance.Api.csproj src/Finance.Api/

RUN dotnet restore src/Finance.Api/Finance.Api.csproj

COPY src ./src

RUN dotnet publish src/Finance.Api/Finance.Api.csproj -c Release -o /out /p:UseAppHost=false

# Runtime stage for the API.
FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app

COPY --from=build /out ./

EXPOSE 8089

ENV ASPNETCORE_URLS=http://+:8089

ENTRYPOINT ["dotnet", "Finance.Api.dll"]
