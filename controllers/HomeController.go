package controllers

import (
	"math"
	"net/url"

	"github.com/astaxie/beego"
	"yy.com/mindoc/conf"
	"yy.com/mindoc/models"
	"yy.com/mindoc/utils/pagination"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	c.Prepare()
	c.TplName = "home/index.tpl"
	//如果没有开启匿名访问，则跳转到登录页面
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(conf.URLFor("AccountController.Login")+"?url="+url.PathEscape(conf.BaseUrl+c.Ctx.Request.URL.RequestURI()), 302)
	}
	pageIndex, _ := c.GetInt("page", 1)
	pageSize := 18

	member_id := 0

	if c.Member != nil {
		member_id = c.Member.MemberId
	}
	books, totalCount, err := models.NewBook().FindForHomeToPager(pageIndex, pageSize, member_id)

	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	if totalCount > 0 {
		pager := pagination.NewPagination(c.Ctx.Request, totalCount, pageSize, c.BaseUrl())
		c.Data["PageHtml"] = pager.HtmlPages()
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	c.Data["Lists"] = books

	labels, totalCount, err := models.NewLabel().FindToPager(1, 10)

	if err != nil {
		c.Data["Labels"] = make([]*models.Label, 0)
	} else {
		c.Data["Labels"] = labels
	}
}
