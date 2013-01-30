

//$(document).ready(function() { formNewUser(); });
$(function(){
    $("#menucategory").change(function(){updateMenu(this.value)});
    updateMenu("user");
    }
);

function updateMenu(name) {
    if (name === "user" && $.mapo.user.isLogedIn()) {
        $("#menucontainer").text("");
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "PROFILE").click(profileUser);
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "LOGOUT").click(userLogout);
    } else if (name === "studio" && $.mapo.user.isLogedIn()) {
        $("#menucontainer").text("");
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "NEW").click(formNewStudio);
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "ALL").click(getStudioList);
    } else if (name === "project" && $.mapo.user.isLogedIn()) {
        $("#menucontainer").text("");
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "NEW").click(formNewProject);
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "ALL").click(getProjectList);
    } else if ($.mapo.user.isLogedIn()) {
        $("#menucontainer").text("");
    } else {
        if (name === "user") {
            $("#menucontainer").text("");
            $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "LOGIN").click(loginUser);
        } else {
            $("#menucontainer").text("");
            $("#menucontainer").text("Please login!");
        }
    }
};

function loginUser() {
    //window.location = "https://accounts.google.com/o/oauth2/auth?scope=https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile&state=profile&redirect_uri=http://localhost:8081/oauth2callback&response_type=code&client_id=60876467348.apps.googleusercontent.com&approval_prompt=force"
    window.location = "https://accounts.google.com/o/oauth2/auth?scope=https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile&state=profile&redirect_uri=http://localhost:8081/oauth2callback&response_type=code&client_id=60876467348.apps.googleusercontent.com"
}

function userLogout() {
    $.cookie("authid", null);
    //updateMenu($("#menucategory")[0].value);
    $.cookie("uid", null);
    $.cookie("sid", null);
    updateMenu($("#menucategory")[0].value);
}

function profileUser() {
    var container = $("#content");
    container.text("");
    $.mapo.widgets.userProfile(container, getJson("/admin/user/"+$.cookie("uid")));
}

function getStudioList() {
    var container = $("#content");
    container.text("");
    $.mapo.widgets.studioList(container, getJson("/admin/studio"));
}

function getStudio(studioid) {
    var dialogcontainer = $("#dialogcontainer");
    var studio = $.mapo.widgets.studio(dialogcontainer, getJson("/admin/studio/"+studioid));
    studio.addClass("formcontainer");
    $.mapo.widgets.errorbox(studio);
    dialogcontainer.dialog({modal: true, width: "660px", close: function(e, u){$("div", this).remove();}});

    $.cookie("sid", studioid);

    //var container = $("#content");
    //container.text("");
}

function formNewStudio() {
    var dialogcontainer = $("#dialogcontainer");
    var form = $.mapo.widgets.form(dialogcontainer, "/admin/studio", "POST");
    form.addClass("formcontainer");

    $.mapo.widgets.title(form, "New Studio");

    $.mapo.widgets.textbox(form, "studioid", "ID:");
    $.mapo.widgets.textbox(form, "name", "Name:");
    $.mapo.widgets.textbox(form, "description", "Description:");

    form.append($("<br/>"));
    form.append($("<br/>"));
    
    $.mapo.widgets.errorbox(form);

    dialogcontainer.dialog({modal: true, width: "660px", close: function(e, u){$("form", this).remove();}});
    $.mapo.widgets.submitbutton(form)
    form.bind("submit", function() {postJson(this); return false;});
}

function formNewProject() {
    var dialogcontainer = $("#dialogcontainer");
    var form = $.mapo.widgets.form(dialogcontainer, "/admin/project", "POST");
    form.addClass("formcontainer");

    $.mapo.widgets.title(form, "New Project");

    $.mapo.widgets.textbox(form, "name", "Name:");
    $.mapo.widgets.textbox(form, "description", "Description:");

    form.append($("<br/>"));
    form.append($("<br/>"));
    
    $.mapo.widgets.errorbox(form);

    dialogcontainer.dialog({modal: true, width: "660px", close: function(e, u){$("form", this).remove();}});
    $.mapo.widgets.submitbutton(form)
    form.bind("submit", function() {postJson(this); return false;});
}

function getProjectList() {
    var container = $("#content");
    container.text("");
    $.mapo.widgets.projectList(container, getJson("/admin/project"));
}

function getProject(projid) {
    var dialogcontainer = $("#dialogcontainer");
    var project = $.mapo.widgets.project(dialogcontainer, getJson("/admin/project/"+projid));
    project.addClass("formcontainer");
    $.mapo.widgets.errorbox(project);
    dialogcontainer.dialog({modal: true, width: "660px", close: function(e, u){$("div", this).remove();}});

    $.cookie("pid", projid);
}

function postJson(form) {
    var errorElements = $('[name^="error"]', form).text(""); 
    //var form = $("form")[0];
    var formdata = new FormData(form);
    var xhr = new XMLHttpRequest();
    xhr.open(form.method, form.action, false);
    //xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    xhr.send(formdata);
    var result = JSON.parse(xhr.responseText);
    var data = result.data;
    if (result.status == "ok") {
        $(form).remove();
        $("#dialogcontainer").dialog("destroy");
        alert("TODO: update page content");
        //updateMenu();
        return false;
    } else {
        for (err in data) {
            var errorDiv = $('[name="error:'+err+'"]', form);
            if (errorDiv.length == 0) {
                $('[name="error"]', form).text(err + ": "+data[err]);
            }
            errorDiv.text(err + ": "+ data[err]);
        };
    };
    return false;
}

function getJson(url) {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, false);
    xhr.send(null);
    var result = JSON.parse(xhr.responseText);
    var data = result.data;
    if (result.status == "ok") {
        return data;
    } else {
        for (err in data) {
            alert("error: "+ data[err]);
        };
    };
    return false;
}
