using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface IUploadJournalEntryImageHandler
{
    Task HandleAsync(int image, UploadJournalEntryImageRequest request, CancellationToken cancellationToken = default);
}
