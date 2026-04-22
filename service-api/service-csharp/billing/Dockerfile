# Build stage for the billing service.
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /src

COPY Directory.Build.props ./
COPY src/Billing.Api/Billing.Api.csproj src/Billing.Api/

RUN dotnet restore src/Billing.Api/Billing.Api.csproj

COPY src ./src

RUN dotnet publish src/Billing.Api/Billing.Api.csproj -c Release -o /out /p:UseAppHost=false

# Runtime stage for the API.
FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app

COPY --from=build /out ./

EXPOSE 8095

ENV ASPNETCORE_URLS=http://+:8095

ENTRYPOINT ["dotnet", "Billing.Api.dll"]
