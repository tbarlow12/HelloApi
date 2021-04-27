using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Web.Mvc;
using HelloApi.Loggers;

namespace HelloApi.Controllers
{
    public class HomeController : Controller
    {
        public ActionResult Index()
        {
            using (var eventLogger = (IEventLogger) HttpContext.GetService(typeof(IEventLogger)))
            {
                eventLogger.LogEvent("SampleEvent", "Someone visited the homepage!");

                ViewBag.Title = "Home Page";

                return View();
            }
        }
    }
}
