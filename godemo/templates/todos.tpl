<div class="row">

    <div class="panel panel-default col-lg-6 col-xs-12">
        <ul class="list-group group-todo">
            <h2>Todo</h2>
            {{range $i, $v := .todo}}
            <li ontouchstart="" class="list-group-item text-wrap">
                <h4 class="list listgroup-item-heading">{{$v.Title}}</h4>
                <p class="list-group-item-text">{{$v.Text}}</p>
                <div class="buttons">
                    <button ontouchstart="" type="button" class="btn btn-info btn-xs" data-toggle="modal" data-target="#exampleModal"
                        title="Edit">
                        <i class="fas fa-edit fa-xs" aria-hidden="true"></i>
                    </button>
                    <button ontouchstart="" type="button" class="btn btn-danger btn-xs move-done" title="Done">
                        <i class="fas fa-check fa-xs class-change" aria-hidden="true"></i>
                    </button>
                </div>
            </li>
            {{end}}
        </ul>
        <div class="panel-footer">
            <small>{{.count}} list items</small>
        </div>
        <!-- Button trigger modal -->
        <button href="#" class="btn btn-success show-todolist-modal" data-toggle="modal" data-target="#exampleModal">Create New List</button>
    </div>

    <div class="panel panel-default col-lg-6 col-xs-12">
            <ul class="list-group group-done text-wrap">
                <h2>Done</h2>
                {{range $i, $v := .done}}
                <li ontouchstart="" class="list-group-item">
                    <h4 class="list listgroup-item-heading">{{$v.Title}}</h4>
                    <p class="list-group-item-text">{{$v.Text}}</p>
                    <div class="buttons">
                        <button ontouchstart="" type="button" class="btn btn-info btn-xs" data-toggle="modal" data-target="#exampleModal"
                            title="Edit">
                            <i class="fas fa-edit fa-xs" aria-hidden="true"></i>
                        </button>
                        <button ontouchstart="" type="button" class="btn btn-danger btn-xs delete-done" title="Delete">
                            <i class="fas fa-minus fa-xs" aria-hidden="true"></i>
                        </button>
                    </div>
                </li>
                {{end}}
            </ul>
            <div class="panel-footer">
                <small>{{.count_done}} list items</small>
            </div>
        </div>
</div>

<!-- Modal -->
<div class="modal fade" id="todolist-modal" tabindex="-1" role="dialog" aria-labelledby="exampleModalCenterTitle"
    aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="exampleModalLongTitle">Create todo list</h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="modal-body">
                <form action="">
                    <div class="form-group">
                        <label for="" class="control-label">List Name</label>
                        <input type="text" class="form-control input-lg input-title">
                    </div>
                    <div class="form-group">
                        <label for="" class="control-label">Description</label>
                        <textarea rows="2" class="form-control"></textarea>
                    </div>
                </form>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                <button type="button" class="btn btn-primary save-change">Save</button>
            </div>
        </div>
    </div>
</div>