using System;
using System.Web.Mvc;
using HelloApi.Loggers;
using StatsdClient;

namespace HelloApi.Filters
{
    public class EventFilter : ActionFilterAttribute, IActionFilter
    {
        public EventFilter() { }

        void IActionFilter.OnActionExecuting(ActionExecutingContext filterContext)
        {
            using (var eventLogger = (IEventLogger) filterContext.HttpContext.GetService(typeof(IEventLogger))) {
                eventLogger.LogEvent(filterContext.RouteData.Route.ToString(), filterContext.HttpContext.Request.ToString());
            }
        }
    }
}
