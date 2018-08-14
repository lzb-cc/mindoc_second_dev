<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>搜索 - Powered by MinDoc</title>

    <!-- Bootstrap -->
    <link href="{{cdncss "/static/bootstrap/css/bootstrap.min.css"}}" rel="stylesheet">
    <link href="{{cdncss "/static/font-awesome/css/font-awesome.min.css"}}" rel="stylesheet">

    <link href="{{cdncss "/static/css/main.css?_=?_=1531986418"}}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="/static/html5shiv/3.7.3/html5shiv.min.js"></script>
    <script src="/static/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>
<div class="manual-reader manual-container manual-search-reader">
    {{template "widgets/header.tpl" .}}
    <div class="container manual-body">
        <div class="search-head">
            <strong class="search-title">显示"{{.Keyword}}"的搜索结果</strong>
        </div>
        <div class="row">
            <div class="manual-list">
                {{range $index,$item := .Lists}}
                <div class="search-item">
                    <div class="title"><a href="{{urlfor "DocumentController.Read" ":key" $item.BookIdentify ":id" $item.Identify}}" title="{{$item.DocumentName}}" >{{str2html $item.DocumentName}}</a> </div>
                    <div class="description">
                        {{str2html $item.Description}}
                    </div>
                    <div class="source">
                        <span class="item">来自：<a href="{{urlfor "DocumentController.Index" ":key" $item.BookIdentify}}" >{{$item.BookName}}</a></span>
                        <span class="item">作者：{{$item.Author}}</span>
                        <span class="item">更新时间：{{date_format  $item.ModifyTime "2006-01-02 15:04:05"}}</span>
                    </div>
                </div>
                {{else}}
                <div class="search-empty">
                    <img src="{{cdnimg "/static/images/search_empty.png"}}" class="empty-image">
					<span class="empty-text">暂无相关搜索结果</span>
                </div>
                {{end}}
                <nav class="pagination-container">
                    {{.PageHtml}}
                </nav>
                <div class="clearfix"></div>
            </div>
        </div>
    </div>
    {{template "widgets/footer.tpl" .}}
</div>
<script src="{{cdnjs "/static/jquery/1.12.4/jquery.min.js"}}"></script>
<script src="{{cdnjs "/static/bootstrap/js/bootstrap.min.js"}}"></script>
</body>
</html>