using System.Web;
using System.Web.Mvc;
using HelloApi.Filters;
using HelloApi.Loggers;
using StatsdClient;

namespace HelloApi
{
    public class FilterConfig
    {
        public static void RegisterGlobalFilters(GlobalFilterCollection filters)
        {
            var statsdConfig = new StatsdConfig()
	        {
                StatsdServerName = "127.0.0.1",
                StatsdPort = 8126,
	        };
            var eventLogger = new DataDogStatsdEventLogger(statsdConfig);
            var logFilter = new LogFilter(eventLogger);
            filters.Add(logFilter);
            filters.Add(new HandleErrorAttribute());
        }
    }
}
