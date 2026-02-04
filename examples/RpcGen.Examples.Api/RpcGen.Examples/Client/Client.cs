using Microsoft.AspNetCore.Mvc;
using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Models;

using System.Net.Http;
using System.Net.Http.Json;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.WebUtilities;

namespace RpcGen.Examples.Endpoints.Client;

public interface IProtocolClient
{
    Task<bool> ChangePasswordAsync(
        string newPassword,
        string oldPassword,
        CancellationToken cancellationToken = default);

    Task<Todo> GetTodoAsync(
        int id,
        CancellationToken cancellationToken = default);

    Task<Todo[]> GetTodosAsync(
        CancellationToken cancellationToken = default);

    Task LoginAsync(
        string password,
        string username,
        CancellationToken cancellationToken = default);

}

public sealed class ProtocolClient : IProtocolClient
{
    private readonly HttpClient _http;

    public ProtocolClient(HttpClient http)
    {
        _http = http;
    }


    public async Task<bool> ChangePasswordAsync(
        string newPassword,
        string oldPassword,
        CancellationToken cancellationToken = default)
    {
        var payload = new ChangePasswordParameters
        {
            NewPassword = newPassword,
            OldPassword = oldPassword,
        };

        var response = await _http.PostAsJsonAsync(
            "/change_password",
            payload,
            cancellationToken);

        response.EnsureSuccessStatusCode();
        var result = await response.Content.ReadFromJsonAsync<bool>(
            cancellationToken: cancellationToken);

        return result!;
    }


    public async Task<Todo> GetTodoAsync(
        int id,
        CancellationToken cancellationToken = default)
    {
        var query = new Dictionary<string, string?>()
        {
            ["id"] = id.ToString(),
        };

        var uri = QueryHelpers.AddQueryString("/get_todo", query);

        var response = await _http.GetAsync(uri, cancellationToken);

        response.EnsureSuccessStatusCode();
        var result = await response.Content.ReadFromJsonAsync<Todo>(
            cancellationToken: cancellationToken);

        return result!;
    }


    public async Task<Todo[]> GetTodosAsync(
        CancellationToken cancellationToken = default)
    {
        var uri = "/get_todos";

        var response = await _http.GetAsync(uri, cancellationToken);

        response.EnsureSuccessStatusCode();
        var result = await response.Content.ReadFromJsonAsync<Todo[]>(
            cancellationToken: cancellationToken);

        return result!;
    }


    public async Task LoginAsync(
        string password,
        string username,
        CancellationToken cancellationToken = default)
    {
        var payload = new LoginParameters
        {
            Password = password,
            Username = username,
        };

        var response = await _http.PostAsJsonAsync(
            "/login",
            payload,
            cancellationToken);

        response.EnsureSuccessStatusCode();
    }


}
