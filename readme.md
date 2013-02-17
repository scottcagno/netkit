netkit
=====
<pre>
    <code>import nk "netkit"</code>
</pre>

web server
---
<pre>
    <code>
        var srv = nk.NewWebServer()
        ...
        srv.Serve(":8080", mux)
    </code>
</pre>

http multiplexer
---
<pre>
    <code>
        var mux = nk.NewServeMux()
        ...
        mux.RouteFunc("GET", "/index", index)
    </code>
</pre>

template loader
---
<pre>
    <code>
        var tmp = nk.NewTemplate("templateFolder", "baseTemplate.html")
        ...
        func init() {
            tmp.LoadTemplates("index.html", etc...)
        }
        ...
        tmp.Render(w, "home.html")
    </code>
</pre>

session store
---
<pre>
    <code>
        var ses = nk.NewSessionStore()
        ...
        s := ses.GetSession(w, r)
        ...
        s.Set("id", &user)
    </code>
</pre>

utilities
---
<pre>
    <code>
        hex := nk.EncodeHex("example")
        str := nk.DecodeHex("6578616d706c65")
        ...
        json := nk.EncodeJSON(userObj)
        nk.DecodeJSON(jsonStr, &userObj)
    </code>
</pre>

mongo data wrapper
---
<pre>
    <code>
        import "labix.org/v2/mgo/bson"
        var dat = nk.NewDataWrapper("127.0.0.1").SetDb("database")
        ...
        dat.SetC("collection")
        ...
        dat.Insert(Object{ attr1, attr2, etc...})
        ...
        dat.Update(bson.M{"_id": obj.Id}, bson.M{"$set", bson.M{ "attr1": attr1, "attr2": attr2, etc... },})
        ...
        dat.Return(1, bson.M{"_id": obj.Id}, &object)
        ...
        dat.Delete(Object)
    </code>
</pre>
