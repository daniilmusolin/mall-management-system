namespace MallClient.Interfaces;

public interface ICrudServiceFactory {
    IApiCrudService<T> GetService<T>() where T : class;
}