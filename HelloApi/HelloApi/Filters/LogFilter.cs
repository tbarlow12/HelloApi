using System;
using System.Web.Mvc;
using HelloApi.Loggers;
using StatsdClient;

namespace HelloApi.Filters
{
    public class LogFilter : IActionFilter
    {
        private readonly IEventLogger _eventLogger;

        public LogFilter(IEventLogger eventLogger) {
            _eventLogger = eventLogger;
    	}

        public void OnActionExecuted(ActionExecutedContext filterContext)
        {
            return; // noop for now
        }

        void IActionFilter.OnActionExecuting(ActionExecutingContext filterContext)
        {
            _eventLogger.LogEvent(filterContext.RouteData.Route.ToString(), filterContext.HttpContext.Request.ToString());
        }
    }
}
