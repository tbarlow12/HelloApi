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
        private readonly IEventLogger logger;

        public HomeController(IEventLogger logger) {
            this.logger = logger;
        }

        public ActionResult Index()
        {
            logger.LogEvent("HomePageEvent", "Here's the home page!");
            ViewBag.Title = "Home Page";

            return View();
        }
    }
}
