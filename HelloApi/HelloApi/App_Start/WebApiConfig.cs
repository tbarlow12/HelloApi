using System;
using System.Collections.Generic;
using System.Linq;
using System.Web.Http;
using HelloApi.Loggers;
using StatsdClient;

namespace HelloApi
{
    public static class WebApiConfig
    {
        public static void Register(HttpConfiguration config)
        {
            // Web API configuration and services

            var statsdConfig = new StatsdConfig() { StatsdServerName = "127.0.0.1" };
            var dogStatsDService = new DataDogStatsdEventLogger(statsdConfig);
            config.Services.Add(typeof(IEventLogger), dogStatsDService);

            // Web API routes
            config.MapHttpAttributeRoutes();

            config.Routes.MapHttpRoute(
                name: "DefaultApi",
                routeTemplate: "api/{controller}/{id}",
                defaults: new { id = RouteParameter.Optional }
            );
        }
    }
}
