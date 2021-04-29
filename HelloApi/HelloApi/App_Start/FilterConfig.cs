using System.Web;
using System.Web.Mvc;
using HelloApi.Filters;
using HelloApi.Loggers;

namespace HelloApi
{
    public class FilterConfig
    {
        public static void RegisterGlobalFilters(GlobalFilterCollection filters)
        {
            var eventLogger = DependencyResolver.Current.GetService<IEventLogger>();
            var logFilter = new LogFilter(eventLogger);

            filters.Add(logFilter);
            filters.Add(new HandleErrorAttribute());
        }
    }
}
