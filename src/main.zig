const std = @import("std");
const web = @import("zhp");

pub const io_mode = .evented;
pub const log_level = .info;

const MainHandler = struct {
    pub fn get(self: *MainHandler, request: *web.Request, response: *web.Response) !void {
        _ = self;
        _ = request;
        try response.headers.put("Content-Type", "text/plain");
        _ = try response.stream.write("Hello, World!");
    }

};

pub const routes = [_]web.Route{
    web.Route.create("home", "/", MainHandler),
};

pub const middleware = [_]web.Middleware{
    web.Middleware.create(web.middleware.LoggingMiddleware),
};

pub fn main() anyerror!void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer std.debug.assert(!gpa.deinit());
    const allocator = gpa.allocator();

    var app = web.Application.init(allocator, .{.debug=true});
    defer app.deinit();

    try app.listen("127.0.0.1", 9000);
    try app.start();
}