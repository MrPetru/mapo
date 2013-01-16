

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
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "UPDATE").click(formUpdateUser);
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "LOGOUT").click(userLogout);
    } else if (name === "studio" && $.mapo.user.isLogedIn()) {
        $("#menucontainer").text("");
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "NEW").click(formNewStudio);
    } else if (name === "project" && $.mapo.user.isLogedIn()) {
        $("#menucontainer").text("");
        $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "NEW").click(formNewProject);
    } else if ($.mapo.user.isLogedIn()) {
        $("#menucontainer").text("");
    } else {
        if (name === "user") {
            $("#menucontainer").text("");
            $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "NEW").click(formNewUser);
            $.mapo.widgets.linkbox($("#menucontainer")[0], "#", "LOGIN").click(formLoginUser);
        } else {
            $("#menucontainer").text("");
            $("#menucontainer").text("Please login!");
        }
    }
};

function formLoginUser() {
    var form = $.mapo.widgets.form($("#dialogcontainer"), "/login", "POST");
    form.addClass("formcontainer");

    $.mapo.widgets.textbox(form, "username", "Username:");
    $.mapo.widgets.textbox(form, "password", "Password:");

    form.append($("<br/>"));
    form.append($("<br/>"));
    $.mapo.widgets.errorbox(form);

    $("#dialogcontainer").dialog({modal: true, width:"660px", close: function(e, u){$("form", this).remove();}});

    $.mapo.widgets.submitbutton(form);
    form.bind("submit", function(){postJson(this); updateMenu($("#menucategory")[0].value); return false;});
}

function userLogout() {
    $.cookie("authid", null);
    updateMenu($("#menucategory")[0].value);
}

function formNewUser() {
    var dialogcontainer = $("#dialogcontainer");
    var form = $.mapo.widgets.form(dialogcontainer, "/admin/user", "POST");
    form.addClass("formcontainer");

    $.mapo.widgets.title(form, "New User Form");

    $.mapo.widgets.textbox(form, "username", "Username:");
    $.mapo.widgets.textbox(form, "password", "Password:");
    $.mapo.widgets.textbox(form, "firstname", "First Name:");
    $.mapo.widgets.textbox(form, "lastname", "Last Name:");
    $.mapo.widgets.textbox(form, "email", "Email:");

    $.mapo.widgets.textareabox(form, "description", "Description:");
    
    form.append($("<br/>"));
    
    $.mapo.widgets.errorbox(form);

    dialogcontainer.dialog({modal: true, width: "660px", close: function(e, u){$("form", this).remove();}});
    $.mapo.widgets.submitbutton(form)
    form.bind("submit", function() {postJson(this); return false;});

    };

function formUpdateUser() {

}

function profileUser() {
    var container = $("#content");
    container.text("");
    $.mapo.widgets.userProfile(container, getJson("/admin/user/"+$.cookie("authid")));
}

function formNewStudio() {
    var dialogcontainer = $("#dialogcontainer");
    var form = $.mapo.widgets.form(dialogcontainer, "/admin/studio", "POST");
    form.addClass("formcontainer");

    $.mapo.widgets.title(form, "New Studio Form");

    $.mapo.widgets.textbox(form, "name", "Studio Name:");

    form.append($("<br/>"));
    form.append($("<br/>"));
    
    $.mapo.widgets.errorbox(form);

    dialogcontainer.dialog({modal: true, width: "660px", close: function(e, u){$("form", this).remove();}});
    $.mapo.widgets.submitbutton(form)
    form.bind("submit", function() {postJson(this); return false;});
}

function formNewProject() {

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
