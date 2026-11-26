namespace MallClient.Interfaces;

public interface IApiErrorParser {
    string Parse(string json);
}