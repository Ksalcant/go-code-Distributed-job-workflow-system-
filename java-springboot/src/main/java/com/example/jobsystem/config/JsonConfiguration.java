package com.example.jobsystem.config;

import com.fasterxml.jackson.databind.cfg.CoercionAction;
import com.fasterxml.jackson.databind.cfg.CoercionInputShape;
import com.fasterxml.jackson.databind.type.LogicalType;
import org.springframework.boot.autoconfigure.jackson.Jackson2ObjectMapperBuilderCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class JsonConfiguration {
    @Bean
    Jackson2ObjectMapperBuilderCustomizer strictStrings() {
        // Go rejects numbers/booleans for string fields; Jackson otherwise converts them.
        return builder -> builder.postConfigurer(mapper -> {
            var strings = mapper.coercionConfigFor(LogicalType.Textual);
            strings.setCoercion(CoercionInputShape.Integer, CoercionAction.Fail);
            strings.setCoercion(CoercionInputShape.Float, CoercionAction.Fail);
            strings.setCoercion(CoercionInputShape.Boolean, CoercionAction.Fail);
        });
    }
}
