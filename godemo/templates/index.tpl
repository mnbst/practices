{{template "header.tpl" .}}
<div class="container">

	<nav class="navbar navbar-inverse navbar-fixed-top">

	  </nav>
	  
	  <!-- title and search -->
	  <div class="section-main">
		  <div class="row">
			<div class="col-xs-12 col-lg-6">
			  <h2>My Todo List</h2>
			</div>
			<div class="col-xs-12 col-lg-6">
			  <form method="GET" class="form-main  form-inline nofloat-xs pull-right pull-left-sm">
	  
				<div class="form-group form-input-fields form-group-lg has-feedback">
				  <label class="sr-only" id="search">Search</label>
				  <div class="input-group">
					<input type="text" class="form-control input-search" name="q" placeholder="Search">
					<span class="input-group-addon group-icon"> 
					</span>
				  </div>
				  <button class="btn btn-md btn-success">
				   <i class="fas fa-search"></i>
				  </button>
				  <button href="/" class="btn btn-md btn-success">See all</button>
				</div>
			  </form>
		  </div>
		</div>
	  </div>

		<a href="/login">Login</a> or <a href="/register">Register</a>
		
		<h2>Usage</h2>
		<ol>
			<li>{{if .user.ID}}<span class="text-success">✓</span>{{end}} <a href="/register">Register</a>, Please enter your email and password.</li>
			<li>{{if .user.ID}}<span class="text-success">✓</span>{{end}} <a href="/login">Login</a>, Enter same email and password.</li>
		</ol>


		{{if .user.ID}}
		{{template "todos.tpl" .}}
		
		<div class="alert alert-warning">
			<dl>
				<dt>Email</dt>
				<dd class="text-wrap">{{.user.Email}}</dd>

				<dt>Password(Hashed)</dt>
				<dd class="text-wrap">{{.user.Password}}</dd>

				<dt>API token</dt>
				<dd class="text-wrap">{{.user.Token}}</dd>

				<dt>CreatedAt</dt>
				<dd class="text-wrap">{{.user.CreatedAt}}</dd>

				<dt>ID</dt>
				<dd class="text-wrap">{{.user.ID}}</dd>
			</dl>

			<a href="/logout" class="btn btn-primary">Log out</a>
			<a href="#" class="btn btn-danger" data-toggle="modal" data-target="#MyModalCenter">Delete registration</a>
		</div>
		{{end}}
</div>

{{template "footer.tpl" .}}

<!-- Modal -->
<div class="modal fade" id="MyModalCenter" tabindex="-1" role="dialog" aria-labelledby="MyModalCenterTitle" aria-hidden="true">
	<div class="modal-dialog modal-dialog-centered" role="document">
	  <div class="modal-content">
		<div class="modal-header">
		  <h5 class="modal-title" id="MyModalLongTitle">Delete your ID & lists</h5>
		  <button type="button" class="close" data-dismiss="modal" aria-label="Close">
			<span aria-hidden="true">&times;</span>
		  </button>
		</div>
		<div class="modal-body">
		  confirmation 
		</div>
		<div class="modal-footer">
		  <button type="button" class="btn btn-secondary" data-dismiss="modal">No</button>
		  <a href="/delregistration" class="btn btn-danger">Yes</a>
		</div>
	  </div>
	</div>
  </div>