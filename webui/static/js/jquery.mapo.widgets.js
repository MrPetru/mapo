
(function() {
    if ($.mapo == undefined) {
        $.mapo = {};
        $.mapo.widgets = {};
        $.mapo.user = {}
    };
})();

var mw = $.mapo.widgets

// create a title
mw.title = function(parent, text) {
    var w = $("<h1/>").text(text).appendTo(parent);
    return w
}

mw.label = function(parent, text, cssClass) {
   if (text == undefined || text == "") {
       text = "None"
   }
   var w = $("<div/>").text(text).appendTo(parent);
   w.addClass(cssClass);

   return w
}

mw.textfield = function(parent, name) {
    var w = $("<input/>", {type:"text", name: name}).appendTo(parent);

    return w
}

mw.submitbutton = function(parent) {
    var w = $("<input/>", {type:"submit"}).appendTo(parent).addClass("submitbutton");

    return w
}

mw.textarea = function(parent, name) {
    var w = $("<textarea/>", {name: name, rows: 5}).appendTo(parent);

    return w
}

mw.errorbox = function(parent, name) {
    if (name != undefined) {
        var w = $("<div/>", {name:"error:" + name}).appendTo(parent).addClass("error");
        return w
    }
    
    var w = $("<div/>", {name:"error"}).appendTo(parent).addClass("error");

    return w
}

mw.textbox = function(parent, name, label) {
    var w = mw.box(parent).addClass("inputBlock");
    mw.label(w, label, "label");
    mw.textfield(w, name).addClass("inputElement");
    mw.errorbox(w, name);

    return w
}

mw.textareabox = function(parent, name, label) {
    var w = mw.box(parent).addClass("inputBlock");
    mw.label(w, label, "label");
    mw.textarea(w, name).addClass("textarea");
    mw.errorbox(w, name);

    return w
}

mw.box = function(parent) {
    var w = $("<div/>").appendTo(parent);

    return w
}

mw.link = function(parent, href, text) {
    var w = $("<a/>", {href:href});
    w.appendTo(parent);
    w.text(text);

    return w
}

mw.linkbox = function(parent, href, text) {
    var w = $("<div/>").addClass("linkbox");
    w.appendTo(parent);
    mw.link(w, href, text);

    return w
}

mw.hidden = function(parent, name, value) {
    var w = $("<input/>", {type: "hidden", name:name, value:value});
    w.appendTo(parent);

    return w
}

mw.form = function(parent, action, method) {
    var w = $("<form/>", {action: action, method: method}).appendTo(parent);

    return w
}

mw.image = function(parent, path, width, height) {
    var w = $("<img/>", {src:path, width:width, height:height}).appendTo(parent);

    return w
}

mw.table = function(parent) {
    var w = $("<table/>").appendTo(parent)

    return w
}

mw.tablerow = function(parent) {
    var w = $("<tr/>").appendTo(parent)

    return w
}

mw.tablecell = function(parent, data) {
    var w = $("<td/>").appendTo(parent)

    if (data != undefined) {
        w.text(data)
    }

    return w
}

mw.userProfile = function(parent, data) {
    var w = mw.box(parent)

    mw.title(w, "Profile Detailes")

    mw.label(w, "ID:", "label");
    mw.label(w, data["Id"], "value")

    mw.label(w, "NAME:", "label");
    mw.label(w, data["Name"], "value")

    mw.label(w, "EMAIL:", "label");
    mw.label(w, data["Email"], "value")

    mw.label(w, "OAUTH PROVIDER:", "label");
    mw.label(w, data["Oauthprovider"], "value")

    mw.label(w, "OAUTH ID:", "label");
    mw.label(w, data["Oauthid"], "value")

    mw.label(w, "AVATAR:", "label");
    mw.image(w, data["picture"], 128, 128)

    return w
}

mw.studio = function(parent, data) {
    var w = mw.box(parent)
    mw.title(w, data["Name"])

    mw.label(w, "ID:", "label")
    mw.label(w, data["Id"])
    mw.label(w, "DESCRIPTION:", "label")
    mw.label(w, data["Description"])
    mw.label(w, "OWNERS ID:", "label")
    mw.label(w, data["Owners"])
    mw.label(w, "PROJECTS:", "label")
    mw.label(w, data["Projects"])

    return w
}

mw.studioList = function(parent, data) {
    var w = mw.table(parent)

    mw.title(w, "List User Studio")

    header = mw.tablerow(w)
    mw.tablecell(header, "ID")
    mw.tablecell(header, "NAME")
    mw.tablecell(header, "DESCRIPTION")
    mw.tablecell(header, "OWNERS")
    mw.tablecell(header, "PROJECTS")

    for (s in data) {
        var row = mw.tablerow(w)
        row.click((function(v) {return function() {getStudio(v)}})(data[s]["Id"]))
        for (c in data[s]) {
            cel = mw.tablecell(row)
            mw.label(cel, data[s][c])
        }
    }

    return w
}

mw.projectList = function(parent, data) {
    var w = mw.table(parent)

    mw.title(w, "List User Studio")

    header = mw.tablerow(w)
    mw.tablecell(header, "ID")
    mw.tablecell(header, "NAME")
    mw.tablecell(header, "DESCRIPTION")
    mw.tablecell(header, "ADMINS")
    mw.tablecell(header, "SUPERVISORS")
    mw.tablecell(header, "ARTISTS")

    for (s in data) {
        var row = mw.tablerow(w)
        row.click((function(v) {return function() {getProject(v)}})(data[s]["Id"]))
        for (c in data[s]) {
            cel = mw.tablecell(row)
            mw.label(cel, data[s][c])
        }
    }

    return w
}

mw.project = function(parent, data) {
    var w = mw.box(parent)
    mw.title(w, data["Name"])

    mw.label(w, "ID:", "label")
    mw.label(w, data["Id"])
    mw.label(w, "DESCRIPTION:", "label")
    mw.label(w, data["Description"])
    mw.label(w, "ADMINS:", "label")
    mw.label(w, data["Admins"])
    mw.label(w, "SUPERVISORS:", "label")
    mw.label(w, data["Supervisors"])
    mw.label(w, "ARTISTS:", "label")
    mw.label(w, data["Artists"])

    return w
}

// helpers

var mu = $.mapo.user

// check if current user is loged in or not
mu.isLogedIn = function() {
    var authid = $.cookie("authid");
    if (authid != null && authid.length > 0) {
        return true;
    }
    return false;
}
