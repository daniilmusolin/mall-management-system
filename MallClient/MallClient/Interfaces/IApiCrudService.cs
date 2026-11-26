namespace MallClient.Interfaces;

public interface IApiCrudService<T> where T : class {
    Task<List<T>> GetAllAsync(int limit = 100, int offset = 0);
    Task CreateAsync(object data);
    Task UpdateAsync(long id, object data);
    Task DeleteAsync(long id);
}