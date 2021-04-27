using System.Web;
using System.Web.Mvc;
using HelloApi.Filters;

namespace HelloApi
{
    public class FilterConfig
    {
        public static void RegisterGlobalFilters(GlobalFilterCollection filters)
        {
            filters.Add(new EventFilter());
            filters.Add(new HandleErrorAttribute());
        }
    }
}
