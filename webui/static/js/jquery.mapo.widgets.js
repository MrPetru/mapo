
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

mw.userProfile = function(parent, data) {
    var w = mw.box(parent)

    mw.title(w, "Profile Detailes")
    mw.label(w, "ID:", "label");
    mw.label(w, data["Id"], "value")

    mw.label(w, "USERNAME:", "label");
    mw.label(w, data["Username"], "value")

    mw.label(w, "FIRST NAME:", "label");
    mw.label(w, data["Firstname"], "value")
    mw.label(w, "LAST NAME:", "label");
    mw.label(w, data["Lastname"], "value")

    mw.label(w, "DESCRIPTION:", "label");
    mw.label(w, data["Description"] || "no description", "value")

    mw.label(w, "RATING:", "label");
    mw.label(w, data["Rating"], "value")

    mw.label(w, "STUDIOS:", "label");
    mw.label(w, data["Studios"], "value")

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
