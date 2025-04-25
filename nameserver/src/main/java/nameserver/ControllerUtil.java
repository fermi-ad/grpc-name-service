package nameserver;
import java.util.ArrayList;
import java.util.List;

import io.quarkus.panache.common.Parameters;

public class ControllerUtil {    
    protected static String createFindParam(List<List<String>> queryParams, Parameters parameters) {            
        List <String> querys = new ArrayList<>();        
        int paramId = 1;
        
        for(var p : queryParams) {
            if (p.size() != 2) {
                throw new RuntimeException("Invalid query parameter: " + p);
            }
            var field = queryParams.get(0);
            var val = queryParams.get(1);
            String paramName = "param" + paramId;
            querys.add("lower(" + field + ") like :" + paramName);
            parameters.and(paramName, "%" + val + "%");
            paramId++;
        }        
        String queryString = String.join(" and ", querys);

        System.out.println("queryString: " + queryString);
        System.out.println("parameters: " + parameters.toString());

        //return all entries if not parameters specified by specifying a condition that's always true
        if (queryString.isEmpty()) {
            queryString = "1=1";
        }
        return queryString;
    }
        
    protected static int[] readPagination(proto.PaginationRequest pagination) {
        int pageSize = 100;
        int pageNum = 1;
        if (pagination != null) {
            if (pagination.getPageSize() > 0) {                
                pageSize = pagination.getPageSize();   
            }
            if (pagination.getPage() > 0) {
                pageNum = pagination.getPage();
            }
        }

        return new int[] {pageNum, pageSize};
    }

    protected static proto.PaginationResponse createPaginationResponse(int[] pagParams, long totalCount) {
        var pageNum = pagParams[0];
        var pageSize = pagParams[1];

        var paginationBuilder = proto.PaginationResponse.newBuilder()
                .setPage(pageNum)
                .setPageSize(pageSize)
                .setTotalCount(totalCount);
        return paginationBuilder.build();
    }
    
}
