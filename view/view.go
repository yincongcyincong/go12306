package view

var ViewHtml = `<!DOCTYPE html>
<html lang="en">
<head>    
	<meta charset="UTF-8">    
	<title>go12306</title>
	<link href="https://cdn.bootcss.com/twitter-bootstrap/2.3.2/css/bootstrap.min.css" rel="stylesheet"><head/>
<body>
		<div>
			<form class="bs-docs-example form-search">
            <div class="input-append">
              始发站
				<select class='span1'>
				  <option>1</option>
				  <option>2</option>
				  <option>3</option>
				  <option>4</option>
				  <option>5</option>
				</select>
            </div>
            <div class="input-prepend">
              终点站
				<select class='span1'>
				  <option>1</option>
				  <option>2</option>
				  <option>3</option>
				  <option>4</option>
				  <option>5</option>
				</select>
            </div>
          </form>
			
			
			日期
			<select class='span1'>
			  <option>1</option>
			  <option>2</option>
			  <option>3</option>
			  <option>4</option>
			  <option>5</option>
			</select>
			<button class="btn" type="button">查询</button>

			乘车人：
			<label class="checkbox inline">
			  <input type="checkbox" id="inlineCheckbox1" value="option1"> 1
			</label>
			<button type="submit" class="btn btn-primary">登陆</button>
		</div>
		<div>
		<table class="table table-bordered">
		  <thead>
			<tr>
			  <th>车次</th>
			  <th>状态</th>
			  <th>始发站</th>
			  <th>终点站</th>
			  <th>出发时间</th>
			  <th>到达时间</th>
			  <th>历时</th>
			  <th>二等座</th>
			  <th>一等座</th>
			  <th>特等座</th>
			  <th>硬座</th>
			  <th>软座</th>
			  <th>硬卧</th>
			  <th>软卧</th>
              <th>操作</th>
			</tr>
		  </thead>
		  <tbody>
			<tr>
			  <td rowspan="2">1</td>
			  <td>Mark</td>
			  <td>Otto</td>
			  <td>@mdo</td>
			</tr>
		  </tbody>
		</table>
		</div>
 	<script src="https://cdn.jsdelivr.net/npm/jquery@1.11.1/dist/jquery.min.js"></script>
	<script src="https://cdn.bootcss.com/twitter-bootstrap/2.3.2/js/bootstrap.min.js"></script>
<body/>
	<script>
		
		

	</script>
<html/>`

